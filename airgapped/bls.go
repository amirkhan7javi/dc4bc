package airgapped

import (
	"encoding/json"
	"fmt"

	"github.com/corestario/kyber/pairing"
	"github.com/corestario/kyber/sign/bls"
	"github.com/corestario/kyber/sign/tbls"
	client "github.com/lidofinance/dc4bc/client/types"
	"github.com/lidofinance/dc4bc/fsm/state_machines/signing_proposal_fsm"
	"github.com/lidofinance/dc4bc/fsm/types/requests"
	"github.com/lidofinance/dc4bc/fsm/types/responses"
)

// handleStateSigningAwaitConfirmations returns a confirmation of participation to create a threshold signature for a data
func (am *Machine) handleStateSigningAwaitConfirmations(o *client.Operation) error {
	var (
		payload responses.SigningProposalParticipantInvitationsResponse
		err     error
	)

	if err = json.Unmarshal(o.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	participantID, err := am.getParticipantID(o.DKGIdentifier)
	if err != nil {
		return fmt.Errorf("failed to get paricipant id: %w", err)
	}
	req := requests.SigningProposalParticipantRequest{
		BatchID:       payload.BatchID,
		ParticipantId: participantID,
		CreatedAt:     o.CreatedAt,
	}
	reqBz, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to generate fsm request: %w", err)
	}

	o.Event = signing_proposal_fsm.EventConfirmSigningConfirmation
	o.ResultMsgs = append(o.ResultMsgs, createMessage(*o, reqBz))
	return nil
}

// handleStateSigningAwaitPartialSigns takes a data to sign as payload and returns a partial sign for the data to broadcast
func (am *Machine) handleStateSigningAwaitPartialSigns(o *client.Operation) error {
	var (
		payload        responses.SigningPartialSignsParticipantInvitationsResponse
		messagesToSign []requests.MessageToSign
		err            error
	)

	if err = json.Unmarshal(o.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err = json.Unmarshal(payload.SrcPayload, &messagesToSign); err != nil {
		return fmt.Errorf("failed to unmarshal messages to sign: %w", err)
	}

	signs := make([]requests.PartialSign, 0, len(messagesToSign))
	participantID, err := am.getParticipantID(o.DKGIdentifier)
	if err != nil {
		return fmt.Errorf("failed to get paricipant id: %w", err)
	}
	for _, m := range messagesToSign {
		partialSign, err := am.createPartialSign(m.Payload, o.DKGIdentifier)
		if err != nil {
			return fmt.Errorf("failed to create partialSign for msg: %w", err)
		}

		signs = append(signs, requests.PartialSign{
			SigningID: m.SigningID,
			Sign:      partialSign,
		})
	}

	req := requests.SigningProposalBatchPartialSignRequests{
		BatchID:       payload.BatchID,
		ParticipantId: participantID,
		PartialSigns:  signs,
		CreatedAt:     o.CreatedAt,
	}

	reqBz, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to generate fsm request: %w", err)
	}

	o.Event = signing_proposal_fsm.EventSigningPartialSignReceived
	o.ResultMsgs = append(o.ResultMsgs, createMessage(*o, reqBz))
	return nil
}

// reconstructThresholdSignature takes broadcasted partial signs from the previous step and reconstructs a full signature
func (am *Machine) reconstructThresholdSignature(o *client.Operation) error {
	var (
		payload         responses.SigningProcessParticipantResponse
		err             error
		messagesPayload []requests.MessageToSign
	)

	if err = json.Unmarshal(o.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	batchPartialSignatures := make(BatchPartialSignatures)
	for _, participant := range payload.Participants {
		for messageID, sign := range participant.PartialSigns {
			batchPartialSignatures.AddPartialSignature(messageID, sign)
		}
	}

	dkgInstance, ok := am.dkgInstances[o.DKGIdentifier]
	if !ok {
		return fmt.Errorf("dkg instance with identifier %s does not exist", o.DKGIdentifier)
	}

	err = json.Unmarshal(payload.SrcPayload, &messagesPayload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal MessagesToSign: %w", err)
	}

	// just convert slice to map
	messages := make(map[string][]byte)
	for _, m := range messagesPayload {
		messages[m.SigningID] = m.Payload
	}
	response := make([]client.ReconstructedSignature, 0, len(batchPartialSignatures))
	for messageID, messagePartialSignatures := range batchPartialSignatures {
		reconstructedSignature, err := am.recoverFullSign(messages[messageID], messagePartialSignatures, dkgInstance.Threshold,
			dkgInstance.N, o.DKGIdentifier)
		if err != nil {
			return fmt.Errorf("failed to reconsruct full signature for msg: %w", err)
		}
		response = append(response, client.ReconstructedSignature{
			SigningID:  messageID,
			Signature:  reconstructedSignature,
			DKGRoundID: o.DKGIdentifier,
			SrcPayload: messages[messageID],
		})
	}

	respBz, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to generate reconstructed signature response: %w", err)
	}
	o.Event = client.SignatureReconstructed
	o.ResultMsgs = append(o.ResultMsgs, createMessage(*o, respBz))
	return nil
}

// createPartialSign returns a partial sign of a given message
// with using of a private part of the reconstructed DKG key of a given DKG round
func (am *Machine) createPartialSign(msg []byte, dkgIdentifier string) ([]byte, error) {
	blsKeyring, err := am.loadBLSKeyring(dkgIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to load blsKeyring: %w", err)
	}

	return tbls.Sign(am.baseSuite.(pairing.Suite), blsKeyring.Share, msg)
}

// recoverFullSign recovers full threshold signature for a message
// with using of a reconstructed public DKG key of a given DKG round
func (am *Machine) recoverFullSign(msg []byte, sigShares [][]byte, t, n int, dkgIdentifier string) ([]byte, error) {
	blsKeyring, err := am.loadBLSKeyring(dkgIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to load blsKeyring: %w", err)
	}

	return tbls.Recover(am.baseSuite.(pairing.Suite), blsKeyring.PubPoly, msg, sigShares, t, n)
}

// VerifySign verifies a signature of a message
func (am *Machine) VerifySign(msg []byte, fullSignature []byte, dkgIdentifier string) error {
	blsKeyring, err := am.loadBLSKeyring(dkgIdentifier)
	if err != nil {
		return fmt.Errorf("failed to load blsKeyring: %w", err)
	}

	return bls.Verify(am.baseSuite.(pairing.Suite), blsKeyring.PubPoly.Commit(), msg, fullSignature)
}

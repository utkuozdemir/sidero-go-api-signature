// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package auth provides client for authentication API.
package auth

import (
	"context"

	"google.golang.org/grpc"

	authpb "github.com/siderolabs/go-api-signature/api/auth"
)

// Client for Management API .
type Client struct {
	conn authpb.AuthServiceClient
}

// NewClient builds a client out of gRPC connection.
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		conn: authpb.NewAuthServiceClient(conn),
	}
}

// RegisterPGPPublicKey registers a PGP public key for the given identity and returns the login URL.
// Registered public key will need to be verified before it can be used for signing.
func (client *Client) RegisterPGPPublicKey(ctx context.Context, email string, publicKey []byte) (string, error) {
	resp, err := client.conn.RegisterPublicKey(ctx, &authpb.RegisterPublicKeyRequest{
		Identity: &authpb.Identity{Email: email},
		PublicKey: &authpb.PublicKey{
			PgpData: publicKey,
		},
	})
	if err != nil {
		return "", err
	}

	return resp.GetLoginUrl(), nil
}

// ConfirmPublicKey confirms a PGP public key for the given identity.
// This endpoint requires a valid JWT token.
func (client *Client) ConfirmPublicKey(ctx context.Context, publicKeyID string) error {
	_, err := client.conn.ConfirmPublicKey(ctx, &authpb.ConfirmPublicKeyRequest{
		PublicKeyId: publicKeyID,
	})

	return err
}

// AwaitPublicKeyConfirmation waits for the public key with the given information to be confirmed for the given email.
func (client *Client) AwaitPublicKeyConfirmation(ctx context.Context, publicKeyID string) error {
	_, err := client.conn.AwaitPublicKeyConfirmation(
		ctx,
		&authpb.AwaitPublicKeyConfirmationRequest{
			PublicKeyId: publicKeyID,
		},
	)

	return err
}

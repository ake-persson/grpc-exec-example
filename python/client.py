#!/usr/bin/env python

from __future__ import print_function

import sys
sys.path.append('../auth')

import grpc

import auth_pb2
import auth_pb2_grpc

def run():
    creds = grpc.ssl_channel_credentials(open('../tls_setup/certs/ca.pem', 'rb').read())
    chan = grpc.secure_channel('localhost:8080', creds)
    stub = auth_pb2_grpc.AuthStub(chan)

    resp = stub.GetPublicKey(auth_pb2.Empty())
    print("Public Key:\n" + "".join(map(chr, resp.pem)))

if __name__ == '__main__':
    run()

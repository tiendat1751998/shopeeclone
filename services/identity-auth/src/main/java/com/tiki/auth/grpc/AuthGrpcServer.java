package com.tiki.auth.grpc;

import com.tiki.auth.service.AuthService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

@Component
public class AuthGrpcServer {

    private static final Logger log = LoggerFactory.getLogger(AuthGrpcServer.class);

    private final AuthService authService;

    public AuthGrpcServer(AuthService authService) {
        this.authService = authService;
    }

    public void validateToken(ValidateTokenRequest request, StreamObserver<ValidateTokenResponse> responseObserver) {
        try {
            boolean valid = authService.validateToken(request.getAccessToken());

            ValidateTokenResponse response = ValidateTokenResponse.newBuilder()
                .setValid(valid)
                .build();

            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("gRPC validate token failed", e);
            responseObserver.onError(
                Status.INTERNAL.withDescription(e.getMessage()).asRuntimeException()
            );
        }
    }

    // DTO classes - in production these would be generated from proto
    public static class ValidateTokenRequest {
        private String accessToken;

        public ValidateTokenRequest(String accessToken) {
            this.accessToken = accessToken;
        }

        public String getAccessToken() {
            return accessToken;
        }
    }

    public static class ValidateTokenResponse {
        private boolean valid;

        public ValidateTokenResponse(boolean valid) {
            this.valid = valid;
        }

        public boolean getValid() {
            return valid;
        }

        public static Builder newBuilder() {
            return new Builder();
        }

        public static class Builder {
            private boolean valid;

            public Builder setValid(boolean valid) {
                this.valid = valid;
                return this;
            }

            public ValidateTokenResponse build() {
                return new ValidateTokenResponse(valid);
            }
        }
    }
}

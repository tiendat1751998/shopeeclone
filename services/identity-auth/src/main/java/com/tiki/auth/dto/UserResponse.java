package com.tiki.auth.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UserResponse {

    private UUID userId;
    private String email;
    private String phone;
    private String fullName;
    private String role;
    private boolean verified;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}

package com.tiki.auth.controller;

import com.tiki.auth.security.JwksProvider;
import lombok.RequiredArgsConstructor;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
public class JwksController {

    private final JwksProvider jwksProvider;

    @GetMapping("/.well-known/jwks.json")
    public ResponseEntity<String> jwks() {
        return ResponseEntity.ok()
            .contentType(MediaType.APPLICATION_JSON)
            .body(jwksProvider.getJwksSetJson());
    }

    @GetMapping("/.well-known/jwk.json")
    public ResponseEntity<String> jwk() {
        return ResponseEntity.ok()
            .contentType(MediaType.APPLICATION_JSON)
            .body(jwksProvider.getJwksJson());
    }
}

package com.shopee.auth.security;

import io.jsonwebtoken.Claims;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

@Component
@RequiredArgsConstructor
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    private static final Logger log = LoggerFactory.getLogger(JwtAuthenticationFilter.class);
    private static final String AUTHORIZATION_HEADER = "Authorization";
    private static final String BEARER_PREFIX = "Bearer ";

    private final JwtTokenProvider jwtTokenProvider;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        try {
            String token = extractToken(request);

            if (StringUtils.hasText(token)) {
                Claims claims = jwtTokenProvider.validateAccessToken(token);
                if (claims != null) {
                    String userId = claims.getSubject();
                    String role = claims.get("role", String.class);
                    String scope = claims.get("scope", String.class);

                    MDC.put("user_id", userId);
                    MDC.put("user_role", role);

                    List<SimpleGrantedAuthority> authorities = buildAuthorities(role, scope);

                    UsernamePasswordAuthenticationToken auth =
                        new UsernamePasswordAuthenticationToken(userId, null, authorities);
                    auth.setDetails(claims);

                    SecurityContextHolder.getContext().setAuthentication(auth);
                }
            }

            filterChain.doFilter(request, response);
        } finally {
            MDC.remove("user_id");
            MDC.remove("user_role");
        }
    }

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        String path = request.getRequestURI();
        return path.startsWith("/actuator") || path.equals("/health") || path.equals("/ready")
            || path.startsWith("/.well-known");
    }

    private String extractToken(HttpServletRequest request) {
        String bearerToken = request.getHeader(AUTHORIZATION_HEADER);
        if (StringUtils.hasText(bearerToken) && bearerToken.startsWith(BEARER_PREFIX)) {
            return bearerToken.substring(BEARER_PREFIX.length());
        }
        return null;
    }

    private List<SimpleGrantedAuthority> buildAuthorities(String role, String scope) {
        List<SimpleGrantedAuthority> authorities = new ArrayList<>();
        authorities.add(new SimpleGrantedAuthority("ROLE_" + role));

        if (StringUtils.hasText(scope)) {
            Arrays.stream(scope.split("\\s+"))
                .map(s -> new SimpleGrantedAuthority("SCOPE_" + s.replace(":", "_")))
                .forEach(authorities::add);
        }

        return authorities;
    }
}

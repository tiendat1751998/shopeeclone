package com.shopee.auth;

import com.shopee.auth.security.JwksProvider;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class JwksProviderTest {

    @Test
    void init_GeneratesRSAKeyPair() {
        JwksProvider provider = new JwksProvider();
        try {
            var field = JwksProvider.class.getDeclaredField("configuredPrivateKey");
            field.setAccessible(true);
            field.set(provider, "test");
            var field2 = JwksProvider.class.getDeclaredField("configuredPublicKey");
            field2.setAccessible(true);
            field2.set(provider, "test");
            provider.init();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        assertThat(provider.getRsaKey()).isNotNull();
        assertThat(provider.getPublicKey()).isNotNull();
        assertThat(provider.getPrivateKey()).isNotNull();
        assertThat(provider.getKeyId()).isNotBlank();
    }

    @Test
    void getJwksJson_ReturnsValidJson() {
        JwksProvider provider = new JwksProvider();
        try {
            var field = JwksProvider.class.getDeclaredField("configuredPrivateKey");
            field.setAccessible(true);
            field.set(provider, "test");
            var field2 = JwksProvider.class.getDeclaredField("configuredPublicKey");
            field2.setAccessible(true);
            field2.set(provider, "test");
            provider.init();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        String json = provider.getJwksJson();
        assertThat(json).isNotBlank();
        assertThat(json).contains("kty");
        assertThat(json).contains("use");
    }

    @Test
    void getJwksSetJson_ReturnsKeySet() {
        JwksProvider provider = new JwksProvider();
        try {
            var field = JwksProvider.class.getDeclaredField("configuredPrivateKey");
            field.setAccessible(true);
            field.set(provider, "test");
            var field2 = JwksProvider.class.getDeclaredField("configuredPublicKey");
            field2.setAccessible(true);
            field2.set(provider, "test");
            provider.init();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        String json = provider.getJwksSetJson();
        assertThat(json).contains("\"keys\":[");
        assertThat(json).contains("kty");
    }
}

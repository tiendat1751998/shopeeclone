package com.tiki.auth;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.tiki.auth.entity.OutboxEvent;
import com.tiki.auth.repository.OutboxEventRepository;
import com.tiki.auth.service.OutboxPublisher;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Captor;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Map;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class OutboxPublisherTest {

    @Mock
    private OutboxEventRepository outboxEventRepository;

    @Captor
    private ArgumentCaptor<OutboxEvent> eventCaptor;

    private OutboxPublisher outboxPublisher;

    @BeforeEach
    void setUp() {
        outboxPublisher = new OutboxPublisher(outboxEventRepository, new ObjectMapper());
    }

    @Test
    void publish_SavesOutboxEvent() {
        Map<String, String> payload = Map.of("user_id", "123", "email", "test@example.com");

        outboxPublisher.publish("user", "123", "user.registered", payload);

        verify(outboxEventRepository).save(eventCaptor.capture());
        OutboxEvent saved = eventCaptor.getValue();

        assertThat(saved.getAggregateType()).isEqualTo("user");
        assertThat(saved.getAggregateId()).isEqualTo("123");
        assertThat(saved.getEventType()).isEqualTo("user.registered");
        assertThat(saved.getProcessed()).isFalse();
        assertThat(saved.getRetryCount()).isZero();
        assertThat(saved.getPayload()).contains("test@example.com");
    }

    @Test
    void markProcessed_UpdatesEvent() {
        OutboxEvent event = OutboxEvent.builder()
            .eventId(UUID.randomUUID())
            .processed(false)
            .retryCount(0)
            .build();

        outboxPublisher.markProcessed(event);

        assertThat(event.getProcessed()).isTrue();
        assertThat(event.getProcessedAt()).isNotNull();
        verify(outboxEventRepository).save(event);
    }

    @Test
    void markFailed_IncrementsRetryCount() {
        OutboxEvent event = OutboxEvent.builder()
            .eventId(UUID.randomUUID())
            .processed(false)
            .retryCount(0)
            .build();

        outboxPublisher.markFailed(event, "Database connection failed");

        assertThat(event.getRetryCount()).isEqualTo(1);
        assertThat(event.getLastError()).isEqualTo("Database connection failed");
        verify(outboxEventRepository).save(event);
    }
}

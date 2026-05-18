package com.shopee.auth.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.shopee.auth.entity.OutboxEvent;
import com.shopee.auth.repository.OutboxEventRepository;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

@Service
@RequiredArgsConstructor
public class OutboxPublisher {

    private static final Logger log = LoggerFactory.getLogger(OutboxPublisher.class);

    private final OutboxEventRepository outboxEventRepository;
    private final ObjectMapper objectMapper;

    @Transactional
    public void publish(String aggregateType, String aggregateId, String eventType, Object payload) {
        try {
            String payloadJson = objectMapper.writeValueAsString(payload);

            OutboxEvent event = OutboxEvent.builder()
                .aggregateType(aggregateType)
                .aggregateId(aggregateId)
                .eventType(eventType)
                .payload(payloadJson)
                .processed(false)
                .retryCount(0)
                .build();

            outboxEventRepository.save(event);

            log.debug("Outbox event published: type={}, aggregateId={}", eventType, aggregateId);
        } catch (JsonProcessingException e) {
            log.error("Failed to serialize outbox event payload", e);
        }
    }

    @Transactional
    public void markProcessed(OutboxEvent event) {
        event.setProcessed(true);
        event.setProcessedAt(LocalDateTime.now());
        outboxEventRepository.save(event);
    }

    @Transactional
    public void markFailed(OutboxEvent event, String error) {
        event.setRetryCount(event.getRetryCount() + 1);
        event.setLastError(error);
        outboxEventRepository.save(event);
    }
}

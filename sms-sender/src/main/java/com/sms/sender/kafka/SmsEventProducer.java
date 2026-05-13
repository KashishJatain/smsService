package com.sms.sender.kafka;
import com.sms.sender.model.SmsEvent;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;
import java.util.concurrent.CompletableFuture;

@Slf4j
@Service
@RequiredArgsConstructor
public class SmsEventProducer{
    private final KafkaTemplate<String, SmsEvent> kafkaTemplate;

    @Value("${kafka.topic.sms-events}")
    private String smsEventsTopic;
    public void publishEvent(SmsEvent event){
        log.info("Publishing SMS event to Kafka: messageId={},userId={},status={}",
                event.getMessageId(), event.getUserId(),event.getStatus());
        CompletableFuture<SendResult<String, SmsEvent>> future =
                kafkaTemplate.send(smsEventsTopic, event.getUserId(), event);
        future.whenComplete((result, ex) -> {
            if (ex != null){
                log.error("Failed to publish SMS event messageId={} to Kafka: {}",
                        event.getMessageId(), ex.getMessage(), ex);
            } else{
                log.info("SMS event published successfully. messageId={}, partition={},offset={}",
                        event.getMessageId(),
                        result.getRecordMetadata().partition(),
                        result.getRecordMetadata().offset());
            }
        });}}
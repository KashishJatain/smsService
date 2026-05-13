package com.sms.sender.config;
import com.sms.sender.model.SmsEvent;
import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.core.ProducerFactory;

@Configuration
public class KafkaConfig{
    @Value("${kafka.topic.sms-events}")
    private String smsEventsTopic;

    @Bean
    public NewTopic smsEventsTopic(){
        return TopicBuilder.name(smsEventsTopic)
                .partitions(3)
                .replicas(1)
                .build();
    }

    @Bean
    public KafkaTemplate<String,SmsEvent> kafkaTemplate(ProducerFactory<String,SmsEvent> producerFactory){
        return new KafkaTemplate<>(producerFactory);
    }}
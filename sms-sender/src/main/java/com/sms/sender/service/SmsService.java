package com.sms.sender.service;
import com.sms.sender.kafka.SmsEventProducer;
import com.sms.sender.model.SmsEvent;
import com.sms.sender.model.SmsRequest;
import com.sms.sender.model.SmsResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import java.util.UUID;

@Slf4j
@Service
@RequiredArgsConstructor
public class SmsService {
    private final BlockListService blockListService;
    private final VendorService vendorService;
    private final SmsEventProducer smsEventProducer;
    public SmsResponse sendSms(SmsRequest request) {
        String messageId = UUID.randomUUID().toString();
        log.info("Processing SMS request: messageId={}, userId={},phoneNumber={}",
                messageId,request.getUserId(),request.getPhoneNumber());


        if(blockListService.isBlocked(request.getUserId())){
            log.warn("SMS blocked for userId={}. Publishing BLOCKED event.",request.getUserId());
            SmsEvent blockedEvent = SmsEvent.of(messageId,request, "BLOCKED","User is in block list");
            smsEventProducer.publishEvent(blockedEvent);
            return SmsResponse.blocked(request.getPhoneNumber());
        }
        VendorService.VendorResult vendorResult = vendorService.sendSms(
                request.getPhoneNumber(), request.getMessage());

        String status = vendorResult.success() ? "SUCCESS" : "FAIL";

        SmsEvent event = SmsEvent.of(messageId, request, status, vendorResult.response());
        smsEventProducer.publishEvent(event);

        if (vendorResult.success()) {
            log.info("SMS sent successfully. messageId={}", messageId);
            return SmsResponse.success(messageId, request.getPhoneNumber());
        } else {
            log.error("Vendor failed to send SMS. messageId={}, reason={}",
                    messageId, vendorResult.response());
            return SmsResponse.failure(vendorResult.response());
        }
    }
}
package com.sms.sender.model;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.time.Instant;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SmsEvent{
    @JsonProperty("messageId")
    private String messageId;
    @JsonProperty("userId")
    private String userId;
    @JsonProperty("phoneNumber")
    private String phoneNumber;
    @JsonProperty("message")
    private String message;
    @JsonProperty("status")
    private String status; 
    @JsonProperty("vendorResponse")
    private String vendorResponse;
    @JsonProperty("timestamp")
    private String timestamp;
    public static SmsEvent of(String messageId,SmsRequest request,String status,String vendorResponse){
        return SmsEvent.builder()
                .messageId(messageId)
                .userId(request.getUserId())
                .phoneNumber(request.getPhoneNumber())
                .message(request.getMessage())
                .status(status)
                .vendorResponse(vendorResponse)
                .timestamp(Instant.now().toString())
                .build();
    }
}
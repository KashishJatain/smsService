package com.sms.sender.model;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SmsResponse{
    @JsonProperty("messageId")
    private String messageId;
    @JsonProperty("status")
    private String status;
    @JsonProperty("message")
    private String message;
    @JsonProperty("phoneNumber")
    private String phoneNumber;
    public static SmsResponse success(String messageId,String phoneNumber){
        return SmsResponse.builder()
                .messageId(messageId)
                .status("SUCCESS")
                .message("SMS sent and queued for storage successfully.")
                .phoneNumber(phoneNumber)
                .build();
    }
    public static SmsResponse blocked(String phoneNumber){
        return SmsResponse.builder()
                .status("BLOCKED")
                .message("User is blocked from receiving SMS.")
                .phoneNumber(phoneNumber)
                .build();
    }
    public static SmsResponse failure(String reason){
        return SmsResponse.builder()
                .status("FAILED")
                .message("Failed to send SMS: " + reason)
                .build();
    }}
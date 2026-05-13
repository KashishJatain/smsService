package com.sms.sender.model;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SmsRequest{
    @NotBlank(message= "phoneNumber is required")
    @Pattern(regexp= "^\\+?[1-9]\\d{1,14}$",message= "Invalid phone number format")
    @JsonProperty("phoneNumber")
    private String phoneNumber;

    @NotBlank(message= "message is required")
    @JsonProperty("message")
    private String message;

    @NotBlank(message= "userId is required")
    @JsonProperty("userId")
    private String userId;
}
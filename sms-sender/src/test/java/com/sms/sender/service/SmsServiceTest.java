package com.sms.sender.service;
import com.sms.sender.kafka.SmsEventProducer;
import com.sms.sender.model.SmsRequest;
import com.sms.sender.model.SmsResponse;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SmsServiceTest{
    @Mock
    private BlockListService blockListService;
    @Mock
    private VendorService vendorService;
    @Mock
    private SmsEventProducer smsEventProducer;

    @InjectMocks
    private SmsService smsService;
    private SmsRequest validRequest;

    @BeforeEach
    void setUp(){
        validRequest = SmsRequest.builder()
                .userId("user-123")
                .phoneNumber("+919876543210")
                .message("Hello, World!")
                .build();
    }

    @Test
    @DisplayName("Should return SUCCESS when user is not blocked and vendor call succeeds")
    void sendSms_success() {
        when(blockListService.isBlocked("user-123")).thenReturn(false);
        when(vendorService.sendSms(anyString(), anyString()))
                .thenReturn(new VendorService.VendorResult(true,"MOCK-REF-001"));

        SmsResponse response= smsService.sendSms(validRequest);

        assertThat(response.getStatus()).isEqualTo("SUCCESS");
        assertThat(response.getMessageId()).isNotNull();
        verify(smsEventProducer, times(1)).publishEvent(any());
    }

    @Test
    @DisplayName("Should return BLOCKED when user is in the block list")
    void sendSms_blockedUser(){
        when(blockListService.isBlocked("user-123")).thenReturn(true);

        SmsResponse response=smsService.sendSms(validRequest);

        assertThat(response.getStatus()).isEqualTo("BLOCKED");
        verify(vendorService, never()).sendSms(anyString(), anyString());
        verify(smsEventProducer,times(1)).publishEvent(
                argThat(event -> "BLOCKED".equals(event.getStatus())));
    }

    @Test
    @DisplayName("Should return FAILED when vendor call fails")
    void sendSms_vendorFailure(){
        when(blockListService.isBlocked("user-123")).thenReturn(false);
        when(vendorService.sendSms(anyString(), anyString()))
                .thenReturn(new VendorService.VendorResult(false, "VENDOR_TIMEOUT"));

        SmsResponse response = smsService.sendSms(validRequest);

        assertThat(response.getStatus()).isEqualTo("FAILED");
        verify(smsEventProducer, times(1)).publishEvent(
                argThat(event -> "FAIL".equals(event.getStatus())));
    }

    @Test
    @DisplayName("Published event should contain all request fields")
    void sendSms_eventContainsCorrectFields() {
        when(blockListService.isBlocked(anyString())).thenReturn(false);
        when(vendorService.sendSms(anyString(), anyString()))
                .thenReturn(new VendorService.VendorResult(true, "MOCK-REF-002"));

        smsService.sendSms(validRequest);

        var captor =ArgumentCaptor.forClass(com.sms.sender.model.SmsEvent.class);
        verify(smsEventProducer).publishEvent(captor.capture());
        var event= captor.getValue();
        assertThat(event.getUserId()).isEqualTo("user-123");
        assertThat(event.getPhoneNumber()).isEqualTo("+919876543210");
        assertThat(event.getMessage()).isEqualTo("Hello, World!");
        assertThat(event.getStatus()).isEqualTo("SUCCESS");
        assertThat(event.getTimestamp()).isNotNull();
    }
}
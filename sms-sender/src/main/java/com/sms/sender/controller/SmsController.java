package com.sms.sender.controller;
import com.sms.sender.model.SmsRequest;
import com.sms.sender.model.SmsResponse;
import com.sms.sender.service.BlockListService;
import com.sms.sender.service.SmsService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@Slf4j
@RestController
@RequestMapping("/v1/sms")
@RequiredArgsConstructor
public class SmsController {
    private final SmsService smsService;
    private final BlockListService blockListService;
    @PostMapping("/send")
    public ResponseEntity<SmsResponse> sendSms(@Valid @RequestBody SmsRequest request) {
        log.info("Received SMS send request for userId={}",request.getUserId());
        SmsResponse response= smsService.sendSms(request);
        HttpStatus status= switch (response.getStatus()) {
            case "SUCCESS" -> HttpStatus.OK;
            case "BLOCKED" -> HttpStatus.FORBIDDEN;
            default -> HttpStatus.BAD_GATEWAY;
        };
        return ResponseEntity.status(status).body(response);
    }
    @PostMapping("/block/{userId}")
    public ResponseEntity<String> blockUser(@PathVariable String userId) {
        blockListService.blockUser(userId);
        return ResponseEntity.ok("userId=" + userId + " has been blocked.");
    }

    @DeleteMapping("/block/{userId}")
    public ResponseEntity<String> unblockUser(@PathVariable String userId) {
        blockListService.unblockUser(userId);
        return ResponseEntity.ok("userId=" + userId + " has been unblocked.");
    }
}
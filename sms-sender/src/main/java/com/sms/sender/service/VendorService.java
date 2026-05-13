package com.sms.sender.service;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import java.util.Random;

@Slf4j
@Service
public class VendorService{
    @Value("${vendor.mock-failure-rate:0.1}")
    private double mockFailureRate;
    private final Random random= new Random();
    public VendorResult sendSms(String phoneNumber,String message){
        log.info("[MOCK VENDOR] Sending SMS to phoneNumber={}", phoneNumber);
        simulateLatency();
        if(random.nextDouble()<mockFailureRate){
            log.warn("[MOCK VENDOR] Simulated failure for phoneNumber={}", phoneNumber);
            return new VendorResult(false, "VENDOR_TIMEOUT: simulated failure");
        }
        String refId= "MOCK-" + System.currentTimeMillis();
        log.info("[MOCK VENDOR] SMS sent successfully. refId={}",refId);
        return new VendorResult(true, refId);
    }
    private void simulateLatency(){
        try{
            Thread.sleep(50 + random.nextInt(100));
        } catch(InterruptedException e){
            Thread.currentThread().interrupt();
        }
    }
    public record VendorResult(boolean success, String response) {}
}
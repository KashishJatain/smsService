package com.sms.sender.service;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class BlockListService{
    private final RedisTemplate<String, String> redisTemplate;
    @Value("${redis.blocked-users-key:blocked_users}")
    private String blockedUsersKey;boolean isBlocked(String userId){
        try{
            Boolean member= redisTemplate.opsForSet().isMember(blockedUsersKey, userId);
            boolean blocked=Boolean.TRUE.equals(member);
            log.debug("Block-list check for userId={}: blocked={}", userId, blocked);
            return blocked;
        } catch (Exception e){
            log.error("Redis unavailable during block-list check for userId={}. Failing open. Error: {}",
                    userId, e.getMessage());
            return false;
        }
    }

    public void blockUser(String userId){
        redisTemplate.opsForSet().add(blockedUsersKey, userId);
        log.info("userId={} added to block list", userId);
    }

    public void unblockUser(String userId){
        redisTemplate.opsForSet().remove(blockedUsersKey, userId);
        log.info("userId={} removed from block list", userId);
    }
}
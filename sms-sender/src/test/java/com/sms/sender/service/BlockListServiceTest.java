package com.sms.sender.service;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.core.SetOperations;
import org.springframework.test.util.ReflectionTestUtils;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class BlockListServiceTest{
    @Mock
    private RedisTemplate<String,String> redisTemplate;
    @Mock
    private SetOperations<String, String> setOperations;

    @InjectMocks
    private BlockListService blockListService;

    @BeforeEach
    void setUp() {
        ReflectionTestUtils.setField(blockListService, "blockedUsersKey","blocked_users");
        when(redisTemplate.opsForSet()).thenReturn(setOperations);
    }
    @Test
    @DisplayName("isBlocked returns true when userId is in the set")
    void isBlocked_whenUserBlocked_returnsTrue(){
        when(setOperations.isMember("blocked_users", "user-bad")).thenReturn(true);
        assertThat(blockListService.isBlocked("user-bad")).isTrue();
    }

    @Test
    @DisplayName("isBlocked returns false when userId is not in the set")
    void isBlocked_whenUserNotBlocked_returnsFalse(){
        when(setOperations.isMember("blocked_users", "user-good")).thenReturn(false);
        assertThat(blockListService.isBlocked("user-good")).isFalse();
    }

    @Test
    @DisplayName("isBlocked fails open (returns false) when Redis is unavailable")
    void isBlocked_redisDown_failsOpen() {
        when(setOperations.isMember(anyString(), anyString()))
                .thenThrow(new RuntimeException("Connection refused"));
        // Should not throw; should fail open
        assertThat(blockListService.isBlocked("user-any")).isFalse();
    }

    @Test
    @DisplayName("blockUser adds userId to the Redis set")
    void blockUser_addsToSet() {
        blockListService.blockUser("user-xyz");
        verify(setOperations).add("blocked_users","user-xyz");
    }

    @Test
    @DisplayName("unblockUser removes userId from the Redis set")
    void unblockUser_removesFromSet() {
        blockListService.unblockUser("user-xyz");
        verify(setOperations).remove("blocked_users","user-xyz");
    }
}
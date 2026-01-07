<script lang="ts" setup>
import { Moon, Sunny } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
// import { useRouter } from 'vue-router'
import { useStore } from '@/store'

const { t } = useI18n()
const store = useStore()
// const router = useRouter()

const isDark = useDark()
const toggleDark = useToggle(isDark)

function setLocale(locale: string) {
  store.setLocale(locale as 'zh' | 'en')
}

function goAccount() {
  // router.push('/account')
  window.open(`/user-center/basic`, '_blank')
}

async function logout() {
  await store.logout()
  window.location.href = `/login`
}
</script>

<template>
  <div class="pr-15px">
    <el-dropdown trigger="click">
      <el-button link ml-10px>
        <el-icon :size="22" :color="isDark ? '#fff' : '#262f3e'">
          <i-mdi-translate />
        </el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="setLocale('zh')">
            简体中文
          </el-dropdown-item>
          <el-dropdown-item @click="setLocale('en')">
            English
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <el-button link ml-10px @click="toggleDark()">
      <el-icon :size="22" :color="isDark ? '#fff' : '#262f3e'">
        <Moon v-if="isDark" />
        <Sunny v-else />
      </el-icon>
    </el-button>
    <el-dropdown trigger="click" style="margin-left: 10px">
      <!-- <el-avatar size="small">
          <el-icon><i-ep-user-filled /></el-icon>
        </el-avatar> -->
      <el-button link>
        {{ store.user?.name }}
        <el-icon :size="22" class="el-icon--right" :color="isDark ? '#fff' : '#262f3e'">
          <i-ep-arrow-down />
        </el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="goAccount">
            <el-icon>
              <i-ep-user />
            </el-icon>
            <span>{{ t('account_center') }}</span>
          </el-dropdown-item>
          <el-dropdown-item @click="logout">
            <el-icon> <i-ep-switch-button /> </el-icon>{{ t('log_out') }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<style lang="scss" scoped></style>

<i18n lang="yaml">
zh:
  account_center: '账户中心'
  log_out: '退出登录'
en:
  account_center: 'Account Center'
  log_out: 'Logout'
</i18n>

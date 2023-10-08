<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { MessageSchema } from './i18n/type';

const { t, locale } = useI18n<{
      message: MessageSchema,
    }, 'en'|'zh-CN'>({ useScope: 'global' })
//console.log(i18n)
//locale.value="zh-CN"

</script>

<template>
  <header class="navbar shadow" id="theNav">
    <span class="nav-item h3 text-center">
      {{ t("z标题") }}
    </span>
    <div class="d-flex align-content-center">
      <div class="input-group mx-2">
        <div class="input-group-prepend">
          <span class="input-group-text">Language:</span>
        </div>

        <select v-model="locale" class="custom-select">
          <option :key="'zh-CN'" :value="'zh-CN'">简体中文</option>
          <option :key="'en'" :value="'en'">English</option>
        </select>
      </div>
      <a href="./Redistrbution.html" class="mx-2 my-auto text-center text-nowrap">{{ t('z许可') }}</a>
    </div>
  </header>

  <main class="d-flex flex-column h-100 mx-auto">
    <RouterView :to="{ name: 'auth' }" v-slot="{ Component }">
      <transition name="v">
        <component :is="Component" />
      </transition>
    </RouterView>


  </main>
  <footer class="footer bg-light py-3 my-footer">
    <div class="container flex-row d-flex px-3" style="justify-content: space-between;">
      <div>
        <span>
          @ aaadddfgh
        </span>
        <a class="ml-2" href="https://github.com/aaadddfgh/clipboard-go-vue/">contrubutors</a>
      </div>
      <a href="https://github.com/aaadddfgh/clipboard-go-vue/">{{ t("z源代码") }}</a>
    </div>
  </footer>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>

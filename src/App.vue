<script setup lang="ts">


</script>

<template>
  <header class="navbar shadow" id="theNav">
    <span class="nav-item h3 text-center">
      剪贴板
    </span>
    <a href="./Redistrbution.html">许可</a>
  </header>

  <main class="d-flex flex-column h-100 mx-auto">
    <RouterView :to="{ name: 'auth'}" v-slot="{ Component }">
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
        <a class="ml-2" href="">contrubutors</a>
      </div>
      <a href="">source</a>
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

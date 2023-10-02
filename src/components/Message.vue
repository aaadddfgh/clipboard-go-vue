<template>

    <div  class="message card text-center shadow rounded-lg bg-secondary text-white">
        <div class="details">
            <slot></slot>
        </div>
    </div>

</template>
<script lang="ts" setup>

</script>
<style scoped>
.message{
    position: absolute;
    width: 8rem;
    top: 5rem;
    left: calc(50vw - 4rem);
    min-height: 3rem;
    display: flex;
    align-content: center;
    justify-content: center;
}

</style>
  
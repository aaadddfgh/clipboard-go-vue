<script setup lang="ts">
import { postContent, getContent } from '@/api';
import { transpotAES } from '@/lib/crypto';
import { useRSAStore } from '@/stores/key';
import { ref } from 'vue';
import {  useRouter } from 'vue-router';
import MessageVue from './Message.vue';
const router = useRouter()
if (useRSAStore().key === "") {
    // router.push({name:"auth"})
}

function postBoard(content: string) {
    try {
        postContent(
            {
                test: transpotAES.encypt("test"),
                content: transpotAES.encypt(content)
            },
        ).then(({ ok }) => {
            if (ok) {
                printContent();
                info("放入成功");
            }
            else {
                alert("出错了")
                router.push({ name: "auth" })
            }
        })
    }
    catch {
        alert("出错了")
        router.push({ name: "auth" })
    }
}

const content = ref("")

function printContent(action?:(recivedData?:string )=>void) {
    getContent(
        (data) => {
            content.value = transpotAES.decypt(data.data);
            if(action){
                action(content.value)
            }
        })

}
printContent();

const showInfo = ref(false)
const infoData = ref("")
function info(msg: string) {
    showInfo.value = false;
    showInfo.value = true;
    infoData.value = msg;
    setTimeout(() => {
        showInfo.value = false
    }, 3000);
}

function copyContent() {
    navigator.clipboard.writeText(content.value);
    info("已复制")
}
</script>

<template>
    <div class="container d-flex flex-column mx-2 p-2 main">
        <Transition name="message">
            <MessageVue v-if="showInfo">
                {{ infoData }}
            </MessageVue>
        </Transition>
        <span class="h4 mb-3">向剪贴板中放入内容</span>
        <textarea class="shadow text-area mb-3" v-model="content"></textarea>
        <div class="mx-auto d-flex">
            <button type="button" class="btn btn-success mx-3" v-on:click="copyContent">复制</button>
            <button type="button" class="btn btn-primary mx-3" v-on:click="postBoard(content)">放入内容</button>
            <button type="button" class="btn btn-primary mx-3" v-on:click="(e) => { printContent(()=>{info('已刷新')});  }">刷新</button>
        </div>

    </div>
</template>

<style scoped>
.main {
    max-width: 100vw;
    min-width: 70vw;
    min-height: 70vh;
}

.text-area {
    min-height: 40vh;
}

.message-enter-active,
.message-leave-active {
    transition: all 0.5s ease;
}

.message-enter-from,
.message-leave-to {
    transform: translateY(-50px);
    opacity: 20;
}
</style>

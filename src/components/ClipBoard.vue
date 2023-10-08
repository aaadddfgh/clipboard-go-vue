<script setup lang="ts">
import { postContent, getContent } from '@/api';
import { transpotAES } from '@/lib/crypto';
import { useRSAStore } from '@/stores/key';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import MessageVue from './Message.vue';
import Quill  from "quill";
const router = useRouter()
//@ts-ignore
var quillInstance ;
import { useI18n } from 'vue-i18n'
import type { MessageSchema } from '@/i18n/type';

const { t } = useI18n<{
      message: MessageSchema,
    }>()
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
                info(t("z放入成功信息"));
            }
            else {
                alert(t("z错误提示"))
                router.push({ name: "auth" })
            }
        })
    }
    catch {
        alert(t("z错误提示"))
        router.push({ name: "auth" })
    }
}

//const content = ref("")


function printContent(action?: (recivedData?: string) => void) {

    getContent(
        (data) => {
            try {
                //content.value = transpotAES.decypt(data.data);
                //@ts-ignore
                setClipboaed(transpotAES.decypt(data.data));
                if (action) {
                    //action(content.value)
                    //@ts-ignore
                    action(/*editDiv.innerHTML*/)
                }
            }
            catch {
                if (confirm(t("z出错重进提示"))) {
                    router.push({ 'name': 'auth' })
                }
            }
        })

}
onMounted(() => {
    quillInstance = new Quill('#editor-container', {
        modules: {
            toolbar: [
                ['bold', 'italic', 'underline', 'strike'],        // toggled buttons
                ['blockquote', 'code-block', "image"],

                [{ 'header': 1 }, { 'header': 2 }],               // custom button values
                [{ 'list': 'ordered' }, { 'list': 'bullet' }],
                [{ 'script': 'sub' }, { 'script': 'super' }],      // superscript/subscript
                [{ 'indent': '-1' }, { 'indent': '+1' }],          // outdent/indent
                [{ 'direction': 'rtl' }],                         // text direction

                [{ 'size': ['small', false, 'large', 'huge'] }],  // custom dropdown
                [],

                [{ 'color': [] }, { 'background': [] }],          // dropdown with defaults from theme
                [{ 'font': [] }],
                [{ 'align': [] }],

                ['clean']                                         // remove formatting button
            ]
        },
        placeholder: t('z输入占位'),
        theme: 'snow',  // or 'bubble'
        scrollingContainer:'#editableDiv'
    });
    //window.qqqq=quillInstance;
    
    printContent();
})


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

function getClipboard(): string {
//@ts-ignore
    return JSON.stringify((quillInstance as Quill).getContents());
}

function setClipboaed(v: string) {
    //@ts-ignore
    (quillInstance as Quill).setContents(JSON.parse(v))
}

function copyContent() {
    
    //@ts-ignore
    (quillInstance as Quill).setSelection(0,(quillInstance as Quill).getLength())
    document.execCommand("copy")
    //navigator.clipboard.writeText(getClipboard());
    info(t("z复制后信息"))
}


const ws = new WebSocket("ws://" + location.host + "/ws");
// 连接成功
ws.onopen = () => {
    console.log('WebSocket连接成功')
}
// 接收到消息
ws.onmessage = (e) => {
    printContent()
    info(t('z刷新后信息'))
};

</script>

<template>
    <div class="container d-flex flex-column mx-2 p-2 main">
        <Transition name="message">
            <MessageVue v-if="showInfo">
                {{ infoData }}
            </MessageVue>
        </Transition>
        <span class="h4 mb-3">{{t('z输入占位')}}</span>
        <div class="shadow text-area mb-3 p-0 border border-dark rounded" id="editableDiv">
            <div id="editor-container"></div>
        </div>
        <div class="mx-auto d-flex">
            <button type="button" class="btn btn-success mx-3" v-on:click="copyContent">{{t("z复制按钮")}}</button>
            <button type="button" class="btn btn-primary mx-3" v-on:click="//@ts-ignore
                postBoard(getClipboard())">{{t("z放入剪贴板按钮")}}</button>
            <button type="button" class="btn btn-primary mx-3"
                v-on:click="() => { printContent(() => { info(t('z刷新后信息')) }); }">{{t("z刷新按钮")}}</button>
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
    max-width: 80vw;
    overflow: auto;
}

#editor-container{
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

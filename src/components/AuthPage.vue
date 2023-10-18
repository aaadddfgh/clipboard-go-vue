<script setup lang="ts">
import { getKey, getAuth } from '@/api';
import {  transpotAES, handshakePasswordEncrypt, RSA } from '@/lib/crypto';
import { useRSAStore } from '@/stores/key';
import { onBeforeMount, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n'
import type { MessageSchema } from '@/i18n/type';

const { t,locale } = useI18n<{
      message: MessageSchema,
    }>()
// const skip=false;
// if(skip){

// }
const loading = ref(false)
const router = useRouter()
const password = ref("")
const needPass = ref(false)
// {
//     router.push( {name:"clipboard"})
//     //router.push("clipboard")
// }
const ServerRSAPub = useRSAStore();
onBeforeMount(async () => {
    //TODO 需要根据信息配置i18n
    const data = (await getKey())
    //@ts-ignore
    locale.value=data.lang ? data.lang:"";
    ServerRSAPub.key = data.key;
    needPass.value = data.password;
    if (!data.password) {

        alert(t("z无密码警告"))
        auth();
    }
})


function authProcess(data: { ok: boolean, /*key?: string*/ }) {
    if (data.ok) {
        // const aes = handshakeAESDecrypt(data.key)
        // transpotAES.key = aes.key;
        // transpotAES.iv = aes.iv;
        // AES密钥由自己生成
        router.push({ 'name': 'clipboard' })
    }
    else {
        alert(t("z错误提示"))
        router.push({ 'name': 'auth' })
        location.reload()
    }
}

function auth(pass?: string) {
    loading.value = true;
    //const MyRsa = RSA;
    try {
        if (pass && needPass.value) {
            getAuth(
                {
                    pass: handshakePasswordEncrypt(pass, ServerRSAPub.key),
                    //pubKey: MyRsa.getPubKey(),
                    aes: transpotAES.getKeyForTransport(ServerRSAPub.key)
                },
                authProcess
            )
        }
        else {
            getAuth(
                {
                    //pubKey: MyRsa.getPubKey(),
                    aes:  transpotAES.getKeyForTransport(ServerRSAPub.key)
                },
                authProcess
            )
        }
    }
    catch {

    }
}
</script>

<template>
    <div class="d-flex text-center m-auto" v-if="loading">
        <div>
            <div class="spinner-border m-auto" role="status">
                <span class="sr-only">Loading...</span>
            </div>
            <h4 class="mx-auto">{{t("z载入中")}}</h4>
        </div>
    </div>

    <div class="card mx-auto" v-else-if="needPass">


        <div class="card-body">
            <h4 class="card-title text-center">{{t("z输入密码")}}</h4>

            <input type="password" v-model="password" />
            <button v-on:click="auth(password)" type="button" class="mr-3 mt-3 btn btn-primary btn-lg btn-block">接入</button>

        </div>


    </div>
    <div class="d-flex text-center m-auto" v-else>
        <div>
            <div class="spinner-border m-auto" role="status">
                <span class="sr-only">Loading...</span>
            </div>
            <h4 class="mx-auto">{{t("z载入中")}}</h4>
        </div>
    </div>
</template>

<style scoped></style>
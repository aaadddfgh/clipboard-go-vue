import { createRouter,createWebHashHistory } from 'vue-router'
import AuthPage from "@/components/AuthPage.vue";
import ClipBoard from "@/components/ClipBoard.vue";

const router = createRouter({
  history:createWebHashHistory(),
  routes: [
    {
      path: '/',
      name:"auth",
      component: AuthPage
    },
    {
      path: '/clipboard',
      name:"clipboard",
      component: ClipBoard,

    },
  ]
})

export default router

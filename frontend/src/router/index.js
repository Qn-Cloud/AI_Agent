import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Chat from '../views/Chat.vue'
import History from '../views/History.vue'
import Settings from '../views/Settings.vue'
import ApiTest from '../views/ApiTest.vue'
import CharacterEditor from '../components/CharacterEditor.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: {
        title: '角色选择'
      }
    },
    {
      path: '/chat',
      name: 'chat',
      component: Chat,
      meta: {
        title: '对话聊天'
      }
    },
    {
      path: '/chat/:characterId',
      name: 'chat-character',
      component: Chat,
      meta: {
        title: '角色对话'
      }
    },
    {
      path: '/history',
      name: 'history',
      component: History,
      meta: {
        title: '对话历史'
      }
    },
    {
      path: '/settings',
      name: 'settings',
      component: Settings,
      meta: {
        title: '角色设置'
      }
    },
    {
      path: '/api-test',
      name: 'api-test',
      component: ApiTest,
      meta: {
        title: 'API测试'
      }
    },
    {
      path: '/api-test',
      name: 'api-test',
      component: ApiTest,
      meta: {
        title: 'API测试'
      }
    },
    {
      path: '/character-editor',
      name: 'character-editor',
      component: CharacterEditor,
      meta: {
        title: '角色编辑'
      }
    }
  ]
})

router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - AI角色对话`
  }
  next()
})

export default router

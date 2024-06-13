import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import UploadPhotoView from '../views/UploadPhotoView.vue'
import BannedListView from '../views/BannedListView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/session', component: HomeView},
		{path: '/users/:idUser', component: ProfileView},
		{path: '/users/:idUser/images/', component: UploadPhotoView},
		{path: '/users/:idUser/bans/', component: BannedListView},
	]
})

export default router

import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import UploadPhotoView from '../views/UploadPhotoView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/session', component: HomeView},
		{path: '/users/:idUser', component: ProfileView},
		{path: '/users/:idUser/images/', component: UploadPhotoView},
	]
})

export default router

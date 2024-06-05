<script>
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),

            profile: {
				idUser: localStorage.token,
				username: localStorage.username,
				followersCount: 0,
				followingCount: 0,
				images: [
					{
						idImage: 0,
						idOwner: 0,
						dateTime: null,
						file: null,
					}
				],
			},
		}
	},
	methods: {
		async doLogout() {
			try {
				localStorage.removeItem("token")
				localStorage.removeItem("username")
				this.$router.push({ path: '/' });
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		async getUserProfile() {
			try {
            	let response = await this.$axios.get("/users/" + this.profile.idUser, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				this.profile.followersCount = response.data.followerCount
				this.profile.followingCount = response.data.followCount
				this.profile.images = response.data.images
				
				for (let i = 0; i < this.profile.images.length; i++) {
					this.profile.images[i].file = 'data:image/*;base64,' + this.profile.images[i].file
				}
			}
			catch(e) {
				this.errormsg = e.toString();
			}
        },
		
	},
	mounted() {
		this.getUserProfile()
	}
}
</script>

<template>
	<div>
		<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
			<div class="position-sticky pt-3 sidebar-sticky">
				<h6
					class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
					<span>General</span>
				</h6>
				<ul class="nav flex-column">
					<li class="nav-item">
						<RouterLink to="/session" class="nav-link">
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#home" />
							</svg>
							Home
						</RouterLink>
					</li>
					<li class="nav-item">
						<RouterLink :to="'/users/' + profile.idUser + '/images/'" class="nav-link">
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#home" />
							</svg>
							Prova Upload
						</RouterLink>
					</li>
				</ul>
			</div>
		</nav>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h3>{{this.username }}'s profile</h3>
            <h5>Follower {{ this.profile.followersCount }}</h5>
            <h5>Following {{ this.profile.followingCount }}</h5>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button class="btn btn-danger" type="button" @click="doLogout">Logout</button>
					<input type="file" accept="image/*" class="btn btn-outline-primary" @change="uploadFile" ref="file">
					<button class="btn btn-success" @click="submitFile">Upload</button>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-md-4" v-for="image in this.profile.images" :key="image.file">
				<!-- {{ image }} -->
				<div class="card mb-4 shadow-sm">
                	<img class="card-img-top" :src=image.file alt="Card image cap">
				</div>
			</div>
		</div>
	</div>
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>


{
	"user": { "idUser": 5, "username": "paper", "biography": "" },
	"followCount": 0,
	"followerCount": 0,
	"idImage": [ 1, 2 ]
}
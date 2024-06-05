<script>
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),
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

		async viewProfile() {
			try {
				this.$router.push({path: 'users/' + this.token})
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},
		
	},
	mounted() {
		
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
				</ul>
			</div>
		</nav>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Welcome back {{this.username }}</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button class="btn btn-danger" type="button" @click="doLogout">Logout</button>
					<button class="btn btn-primary" type="button" @click="viewProfile">Profile</button>
					<input type="file" accept="image/*" class="btn btn-outline-primary" @change="uploadFile" ref="file">
					<button class="btn btn-success" @click="submitFile">Upload</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
</style>

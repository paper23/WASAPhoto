<script>
import Navbar from '@/components/Navbar.vue';
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),

			showModal: false,
			commentText: "",
			tmpIdImageModal: null,
			tmpIdImageDropDown: null,
			showDropDown: false,

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
						likesCount: 0,
						commentsCount: 0,
						likeStatus: null,
						comments: [
							{
								commentData: [
									{
										idComment: 0,
										idUserWriter: 0,
										idImage: 0,
										text: "",
									}
								],
								username: "",
								
							}
						],
					}
				],
				photoCount: 0,
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

				this.profile.photoCount = this.profile.images.length
			}
			catch(e) {
				this.errormsg = e.toString();
			}
        },
		
		async toggleLike(idImage, likeStatus) {
			if (likeStatus) {
				await this.unlikePhoto(idImage)
			}
			else {
				await this.likePhoto(idImage)
			}

			window.location.reload();
		},

		async likePhoto(idImage) {
			try {
            	let response = await this.$axios.post("/users/" + this.profile.idUser + "/images/" + idImage + "/likes/", {}, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				
				for (let i = 0; i < this.profile.images.length; i++) {
					if (this.profile.images[i].idImage = idImage) {
						this.profile.images[i].likeStatus = !this.profile.images[i].likeStatus
						break;
					}
				}
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		async unlikePhoto(idImage) {
			try {
            	let response = await this.$axios.delete("/users/" + this.profile.idUser + "/images/" + idImage + "/likes/" + this.token, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				
				for (let i = 0; i < this.profile.images.length; i++) {
					if (this.profile.images[i].idImage = idImage) {
						this.profile.images[i].likeStatus = !this.profile.images[i].likeStatus
						break;
					}
				}
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		closeModal() {
      		this.showModal = false;
    	},

		async submitComment(text, idImage) {
			try {
            	let response = await this.$axios.post("/users/" + this.profile.idUser + "/images/" + idImage + "/comments/", {text}, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				this.showModal = false
				window.location.reload()
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},
		
	},
	mounted() {
		this.getUserProfile()
	},
	components: {
		Navbar,
	},
}
</script>

<template>
	<div>
		<Navbar />
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h3>{{this.username }}'s profile</h3>
			<h5>Photos {{ this.profile.photoCount }}</h5>
            <h5>Follower {{ this.profile.followersCount }}</h5>
            <h5>Following {{ this.profile.followingCount }}</h5>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button class="btn btn-danger" type="button" @click="doLogout">Logout</button>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-md-4" v-for="image in this.profile.images" :key="image.idImage">
				<div class="card mb-4 shadow-sm fixed-size">
                	<img class="card-img-top" :src=image.file alt="Card image cap">
				</div>
				<div class="d-flex justify-content-between align-items-center mb-2">
                        <p class="card-text mb-0">Likes : {{ image.likesCount }}</p>
						<button :class="['btn', image.likeStatus ? 'btn-success' : 'btn-danger', 'btn-sm']" type="button" @click="toggleLike(image.idImage, image.likeStatus)">{{ image.likeStatus ? 'Unlike' : 'Like' }}</button>
                </div>
				<div class="d-flex justify-content-between align-items-center mb-2">
                        <p class="card-text mb-0 clickable" @click="showDropDown = !showDropDown; tmpIdImageDropDown = image.idImage">
							Comments : {{ image.commentsCount }}
							<svg class="feather" v-if="showDropDown && tmpIdImageDropDown == image.idImage">
								<use href="/feather-sprite-v4.29.0.svg#chevron-up" />
							</svg>
							<svg class="feather" v-else>
								<use href="/feather-sprite-v4.29.0.svg#chevron-down" />
							</svg>
						</p>
						<button class="btn btn-secondary btn-sm" type="button" @click="showModal = true; tmpIdImageModal = image.idImage">Comment</button>
                </div>
				<div class="d-flex justify-content-between align-items-center mb-2">
					<div class="dropdown">
						<div class="dropdown-content" v-if="showDropDown && tmpIdImageDropDown == image.idImage">
							<ul v-if="image.comments">
									<li v-for="(item, index) in image.comments" :key="index"><b>{{ item.username }}</b> : {{ item.commentData.text }}</li>
							</ul>
						</div>
					</div>
				</div>
				<div v-if="showModal" class="modal-overlay">
					<div class="modal-content">
						<textarea v-model="commentText" placeholder="Enter your comment"></textarea>
						<button class="btn btn-primary btn-sm" @click="submitComment(commentText, tmpIdImageModal)">Submit</button>
						<button class="btn btn-secondary btn-sm" @click="closeModal">Close</button>
					</div>
				</div>
			</div>
		</div>
	</div>

	

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
	.fixed-size {
		width: 70%;
		height: 70%;
		object-fit: contain;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1050;
	}

	.modal-content {
		background: white;
		padding: 50px;
		border-radius: 15px;
		box-shadow: 0 10px 50px rgba(0, 0, 0, 1);
		text-align: center;
		width: 85%;
	}

	textarea {
		width: 100%;
		height: 100px;
		margin-bottom: 10px;
	}

	.clickable {
		cursor: pointer;
		color: blue;
		text-decoration: underline;
	}

	.clickable:hover {
		color: darkblue;
	}

</style>
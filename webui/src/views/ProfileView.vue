<script>
import Navbar from '@/components/Navbar.vue';
import Toolbar from '@/components/Toolbar.vue';
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),

			showModal: false,
			commentText: "",
			tmpIdImageModal: null,
			showDropDown: false,

            profile: {
				idUser: this.$route.params.idUser,
				username: this.$route.query.username,
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
						showDropDownComment: false,
					}
				],
				photoCount: 0,
			},
		}
	},
	methods: {
		async getUserProfile() {
			try {
            	let response = await this.$axios.get("/users/" + this.profile.idUser, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				//this.profile.username = response.data.username
				this.profile.followersCount = response.data.followerCount
				this.profile.followingCount = response.data.followCount
				this.profile.images = response.data.images

				if (this.profile.images == null) {
					return
				}
				
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

		async uncommentPhoto(idImage, idComment) {
			try {
            	let response = await this.$axios.delete("/users/" + this.profile.idUser + "/images/" + idImage + "/comments/" + idComment, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				
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
		Toolbar,
	},
	computed: {
		idUser() {
			return this.$route.params.idUser;
		},
		username() {
		return this.$route.query.username;
		},
  	},
}
</script>

<template>
	<div>
		<Navbar />
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h3>{{this.profile.username }}'s profile</h3>
			<h5>Photos {{ this.profile.photoCount }}</h5>
            <h5>Follower {{ this.profile.followersCount }}</h5>
            <h5>Following {{ this.profile.followingCount }}</h5>
			<Toolbar />
		</div>
		<div class="row" v-if="this.profile.images != null">
			<div class="col-md-4" v-for="image in this.profile.images" :key="image.idImage">
				<div class="card mb-4 shadow-sm fixed-size">
                	<img class="card-img-top" :src=image.file alt="Card image cap">
					<svg v-if="image.idOwner == this.token" class="feather clickable-red position-absolute top-0 end-0 m-0 remove-icon" title="Delete this photo" @click="deletePhoto(image.idImage)">
						<use href="/feather-sprite-v4.29.0.svg#x" />
					</svg>
				</div>
				<div class="d-flex justify-content-between align-items-center mb-2">
                        <p class="card-text mb-0">Likes : {{ image.likesCount }}</p>
						<button :class="['btn', image.likeStatus ? 'btn-success' : 'btn-danger', 'btn-sm']" type="button" @click="toggleLike(image.idImage, image.likeStatus)">{{ image.likeStatus ? 'Unlike' : 'Like' }}</button>
                </div>
				<div class="d-flex justify-content-between align-items-center mb-2">
                        <p class="card-text mb-0 clickable" @click="image.showDropDownComment = !image.showDropDownComment">
							Comments : {{ image.commentsCount }}
							<svg class="feather" v-if="image.showDropDownComment">
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
						<div class="dropdown-content" v-if="image.showDropDownComment">
							<ul v-if="image.comments">
									<li v-for="(item, index) in image.comments" :key="index">
										<b>{{ item.username }}</b> : {{ item.commentData.text }}
										<svg class="feather clickable-red" v-if="item.commentData.idUserWriter == this.token" @click="uncommentPhoto(image.idImage, item.commentData.idComment)">
											<use href="/feather-sprite-v4.29.0.svg#trash-2" />
										</svg>
									</li>
									<br><br><br><br>
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
		<div v-else class="centered">
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#image" />
			</svg>
			<h4>Non hai caricato nessuna immagine</h4>
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
	}

	.clickable:hover {
		color: darkblue;
	}

	.clickable-red {
		cursor: pointer;
		color: red;
	}

	.clickable-red:hover {
		color: darkred;
		text-decoration: underline;
		text-decoration-color: darkred;
	}

	.centered {
		justify-content: center;
		text-align: center;
		align-items: center;
	}

	.remove-icon {
		top: -10px; /* Posiziona l'icona sopra la card */
		right: -10px; /* Posiziona l'icona a destra della card */
		cursor: pointer;
		width: 3px; /* Dimensione dell'icona */
		height: 3px; /* Dimensione dell'icona */
		background: white; /* Sfondo bianco per maggiore contrasto */
		border-radius: 35%; /* Forma circolare */
		padding: 1px; /* Spazio intorno all'icona */
		box-shadow: 0 0 5px rgba(0, 0, 0, 0.4); /* Aggiunge un'ombra per staccarla visivamente */
	}

	.remove-icon use {
		fill: red; /* Colore dell'icona */
	}

</style>
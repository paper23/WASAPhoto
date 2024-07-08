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
			tmpIdUserModal: null,
			showDropDown: false,

			stream: [
				{
					image: {
							idImage: 0,
							idOwner: 0,
							dateTime: null,
							file: null,
							likesCount: 0,
							commentsCount: 0,
							likesCount: 0,
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
						},
					username: "",
					showDropDownComment: false,
				}
			],
		}
	},

	methods: {
		async getMyStream() {
			try {
            	let response = await this.$axios.get("/users/" + this.token + "/stream", {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})

				this.stream = response.data

				if (this.stream == null) {
					return
				}
				
				for (let i = 0; i < this.stream.length; i++) {
					this.stream[i].image.file = 'data:image/*;base64,' + this.stream[i].image.file
				}
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		async toggleLike(idImage, likeStatus, idOwnerImage) {
			if (likeStatus) {
				await this.unlikePhoto(idImage, idOwnerImage)
			}
			else {
				await this.likePhoto(idImage, idOwnerImage)
			}

			window.location.reload();
		},

		async likePhoto(idImage, idOwnerImage) {
			try {
            	let response = await this.$axios.post("/users/" + idOwnerImage + "/images/" + idImage + "/likes/", {}, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		async unlikePhoto(idImage, idOwnerImage) {
			try {
            	let response = await this.$axios.delete("/users/" + idOwnerImage + "/images/" + idImage + "/likes/" + this.token, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		async uncommentPhoto(idImage, idComment, idOwnerImage) {
			try {
            	let response = await this.$axios.delete("/users/" + idOwnerImage + "/images/" + idImage + "/comments/" + idComment, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
				
				window.location.reload()
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		},

		async commentPhoto(text, idImage, idOwnerImage) {
			if (text == "") {
				this.errormsg = "Comment cannot be empty"
			}
			else {
				try {
					let response = await this.$axios.post("/users/" + idOwnerImage + "/images/" + idImage + "/comments/", {text}, {
							headers: {
								Authorization: "Bearer " + localStorage.getItem("token")
							}})
					this.showModal = false
					window.location.reload()
				}
				catch(e) {
					this.errormsg = e.toString();
				}
			}
		},

		formattedDate(dataIn) {
			let date = new Date(dataIn);
			let options = { 
				year: 'numeric', 
				month: 'long', 
				day: 'numeric', 
				hour: 'numeric', 
				minute: 'numeric', 
				second: 'numeric', 
				hour12: true,
			};
			return date.toLocaleString('en-US', options);
    	},
	},

	mounted() {
		this.getMyStream()
	},

	components: {
		Navbar,
		Toolbar,
	},
}
</script>

<template>
	<div>
		<Navbar />
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Welcome back {{this.username }}</h1>
			<Toolbar />
		</div>
		<div class="row" v-if="this.stream != null">
			<div class="col-md-4" v-for="img in this.stream" :key="img.image.idImage">
				<h5>{{ img.username }}</h5>
				<div class="card mb-4 shadow-sm fixed-size">
                	<img class="card-img-top" :src=img.image.file alt="Card image cap">
				</div>
				<div class="d-flex justify-content-between align-items-center mb-2">
					<p class="card-text mb-0">Posted on: {{ formattedDate(img.image.dateTime) }}</p>
				</div>
				<div class="d-flex justify-content-between align-items-center mb-2">
                        <p class="card-text mb-0">Likes : {{ img.image.likesCount }}</p>
						<button :class="['btn', img.image.likeStatus ? 'btn-success' : 'btn-danger', 'btn-sm']" type="button" @click="toggleLike(img.image.idImage, img.image.likeStatus, img.image.idOwner)">{{ img.image.likeStatus ? 'Unlike' : 'Like' }}</button>
                </div>
				<div class="d-flex justify-content-between align-items-center mb-2">
                        <p class="card-text mb-0 clickable" @click="img.showDropDownComment = !img.showDropDownComment">
							Comments : {{ img.image.commentsCount }}
							<svg class="feather" v-if="img.showDropDownComment">
								<use href="/feather-sprite-v4.29.0.svg#chevron-up" />
							</svg>
							<svg class="feather" v-else>
								<use href="/feather-sprite-v4.29.0.svg#chevron-down" />
							</svg>
						</p>
						<button class="btn btn-secondary btn-sm" type="button" @click="showModal = true; tmpIdImageModal = img.image.idImage; tmpIdUserModal = img.image.idOwner">Comment</button>
                </div>
				<div class="d-flex justify-content-between align-items-center mb-2">
					<div class="dropdown">
						<div class="dropdown-content" v-if="img.showDropDownComment">
							<ul v-if="img.image.comments">
									<li v-for="(item, index) in img.image.comments" :key="index">
										<b>{{ item.username }}</b> : {{ item.commentData.text }}
										<svg class="feather clickable-red" v-if="item.commentData.idUserWriter == this.token" @click="uncommentPhoto(img.image.idImage, item.commentData.idComment, img.image.idOwner)">
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
						<button class="btn btn-primary btn-sm" @click="commentPhoto(commentText, tmpIdImageModal, tmpIdUserModal)">Submit</button>
						<button class="btn btn-secondary btn-sm" @click="this.showModal = false">Close</button>
					</div>
				</div>
			</div>
		</div>
		<div v-else class="centered">
			<h4>You still don't follow anyone | None of the users you follow have posted photos</h4>
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

	.centered {
		justify-content: center;
		text-align: center;
		align-items: center;
	}
</style>

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
                },
            }
        }
    }
</script>

<template>
    <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
		<div class="position-sticky pt-3 sidebar-sticky">
			<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
			    <span>Menù</span>
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
							<use href="/feather-sprite-v4.29.0.svg#image" />
						</svg>
						Upload Photo
					</RouterLink>
				</li>
				<li class="nav-item">
					<RouterLink :to="'/users/' + profile.idUser + '/images/'" class="nav-link">
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#search" />
						</svg>
						Search User DA FARE
					</RouterLink>
				</li>
			</ul>
		</div>
	</nav>
</template>
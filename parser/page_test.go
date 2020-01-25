package yp

import (
	"strings"
	"testing"
)

const postsHtml = `
<div class="yp-posts-paginate indigo lighten-5">
	<div class="row">
		<div class="col s12 m5 center">
			<div class="yp-posts-paginate-buttons">
				<a class="btn disabled"><i class="material-icons left">skip_previous</i>First</a>
				<a class="btn disabled""><i class="material-icons left">fast_rewind</i>Previous</a>
			</div>
		</div>
		<div class="col s12 m2 center">
			<div class="center-align">
				<p class="paginate-count">1 / 9</p>
			</div>
		</div>
		<div class="col s12 m5 center">
			<div class="yp-posts-paginate-buttons">
				<a href="/patreon/2533136?p=2" class="btn pag-btn" data-pag="2"><i class="material-icons right">fast_forward</i>Next</a>
				<a href="/patreon/2533136?p=9" class="btn pag-btn" data-pag="9"><i class="material-icons right">skip_next</i>Last</a>
			</div>
		</div>
	</div>
</div>
<div class="row yp-posts-row">
	<div class="col s12 m6">
		<div class="card large yp-post" id="p30853352">
			<div class="card-image waves-effect waves-block waves-light">
				<img class="activator lazyload" src="/static/img/file_placeholder.png" data-src="/patreon_data/2533136/30853352/thumb.jpg" alt="Winter is coming">
			</div>
			<div class="card-content">
				<span class="card-title activator grey-text text-darken-4">Winter is coming<i class="material-icons right">more_vert</i></span>
				<br><span style="font-size: 85%;" class="grey-text post-time">2019-10-18T12:21:13+00:00</span>
			</div>
			<div class="card-action">
				<a href="/patreon_data/2533136/30853352/FLOOF.png" target="_blank">Post file</a>
				<i class="material-icons right yp-flag" title="Flag this post" data-pid="30853352">flag</i>
			</div>
			<div class="card-reveal">
				<span class="card-title grey-text text-darken-4">Winter is coming <small style="font-size: 60%;" class="post-time">2019-10-18T12:21:13+00:00</small><i class="material-icons right">close</i></span>
				<div class="post-body"><p>Yup.</p></div>
				<div class="card-attachments">
					<hr style="height:1px;border:none;color:#212121;background-color:#212121;">
					<span class="card-title grey-text text-darken-4">Media (1)</span>
					<p style="margin-top:-15px;">
					<br><a href="https://yiff.party/patreon_data/2533136/30853352/44578641/FLOOF.png" target="_blank">FLOOF.png</a> (629.2KiB)
					</p>
				</div>
			</div>
		</div>
	</div>
	<div class="col s12 m6">
		<div class="card large yp-post" id="p30306487">
			<div class="card-image waves-effect waves-block waves-light">
				<img class="activator lazyload" src="/static/img/file_placeholder.png" data-src="/patreon_data/2533136/30306487/thumb.jpg" alt="Ninabun, mech pilot">
			</div>
			<div class="card-content">
				<span class="card-title activator grey-text text-darken-4">Ninabun, mech pilot<i class="material-icons right">more_vert</i></span>
				<br><span style="font-size: 85%;" class="grey-text post-time">2019-09-28T08:02:36+00:00</span>
			</div>
			<div class="card-action">
				<a href="/patreon_data/2533136/30306487/ninabunmech.png" target="_blank">Post file</a>
				<i class="material-icons right yp-flag" title="Flag this post" data-pid="30306487">flag</i>
			</div>
			<div class="card-reveal">
				<span class="card-title grey-text text-darken-4">Ninabun, mech pilot <small style="font-size: 60%;" class="post-time">2019-09-28T08:02:36+00:00</small><i class="material-icons right">close</i></span>
				<div class="post-body"><p>Had some hankering of drawing Ninabun again, with some mech. Yeah, heavy inspiration from BL3's Iron Bear. :p Might make it into Iron Bun, though.</p></div>
				<div class="card-attachments">
					<hr style="height:1px;border:none;color:#212121;background-color:#212121;">
					<span class="card-title grey-text text-darken-4">Media (1)</span>
					<p style="margin-top:-15px;">
					<br><a href="https://yiff.party/patreon_data/2533136/30306487/42971270/ninabunmech.png" target="_blank">ninabunmech.png</a> (1021.2KiB)
					</p>
				</div>
				<div class="card-comments">
					<span class="card-title grey-text text-darken-4 yp-post-comment-title">Comments (1)</span>
					<div class="yp-post-comment">
						<div>
							<img class="yp-post-comment-avatar lazyload" data-src="/avatars/patreon/thumb/325476.jpg" src="/static/img/avatar_placeholder.png" alt="user avatar">
						</div>
						<div><span class="yp-post-comment-head"><a href="https://www.patreon.com/user?u=325476" rel="noreferrer" target="_blank"><strong>User #325476</strong></a> - <span class="yp-post-comment-time" data-utc="1569665070">28 Sep 19 10:04</span></span>
							<div class="yp-post-comment-body">Good stuff!</div></div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
<div class="yp-posts-paginate indigo lighten-5">
	<div class="row">
		<div class="col s12 m5 center">
			<div class="yp-posts-paginate-buttons">
				<a class="btn disabled"><i class="material-icons left">skip_previous</i>First</a>
				<a class="btn disabled""><i class="material-icons left">fast_rewind</i>Previous</a>
			</div>
		</div>
		<div class="col s12 m2 center">
			<div class="center-align">
				<p class="paginate-count">1 / 9</p>
			</div>
		</div>
		<div class="col s12 m5 center">
			<div class="yp-posts-paginate-buttons">
				<a href="/patreon/2533136?p=2" class="btn pag-btn pag-btn-bottom" data-pag="2"><i class="material-icons right">fast_forward</i>Next</a>
				<a href="/patreon/2533136?p=9" class="btn pag-btn pag-btn-bottom" data-pag="9"><i class="material-icons right">skip_next</i>Last</a>
			</div>
		</div>
	</div>
</div>

`

func TestParsePosts(t *testing.T) {
	var c = new(Creator)
	posts, err := c.parsePosts(strings.NewReader(postsHtml))
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatal("posts count is whack: ", len(posts))
	}
}

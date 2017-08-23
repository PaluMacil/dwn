import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BlogMode } from '../../services/BlogAPI.service';

@Component({
  selector: 'blog-page',
  templateUrl: './Blog.page.html',
  styleUrls: ['./Blog.page.css']
})
export class BlogPage implements OnInit { 
  mode: BlogMode;
  BlogMode = BlogMode;

  constructor(route: ActivatedRoute, private router: Router) { 
    this.mode = route.snapshot.data.mode;
  }

  ngOnInit() {
  }

  newPost() {
    this.router.navigateByUrl("/blog/new/post");
  }
}
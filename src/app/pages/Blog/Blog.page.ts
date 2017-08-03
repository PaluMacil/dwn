import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'blog-page',
  templateUrl: './Blog.page.html',
  styleUrls: ['./Blog.page.css']
})
export class BlogPage implements OnInit { 
  mode: BlogMode;
  BlogMode = BlogMode;

  constructor(route: ActivatedRoute) { 
    this.mode = route.snapshot.data.mode;
  }

  ngOnInit() {
  }
}

export enum BlogMode {
  SinglePost,
  Topic,
  All
}
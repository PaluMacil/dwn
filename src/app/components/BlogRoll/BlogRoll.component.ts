import { Component, OnInit, Input } from '@angular/core';
import { BlogMode, BlogRoll, Post, BlogAPI } from '../../services/BlogAPI.service';

@Component({
  selector: 'blog-roll',
  templateUrl: './BlogRoll.component.html',
  styleUrls: ['./BlogRoll.component.css']
})
export class BlogRollComponent implements OnInit { 
  @Input('mode') mode: BlogMode;
  public posts = new Array<Post>();

  constructor(public blogAPI: BlogAPI) {}

  ngOnInit() {
    this.blogAPI.GetBlogRoll().subscribe(
      (r)=>{
        this.posts = r.Posts;
      },
      (err)=>{},
      ()=>{},
    );
  }
}
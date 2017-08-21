import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { BlogAPI, Post, PostFormat } from '../../services/BlogAPI.service'

@Component({
  selector: 'post-editor',
  templateUrl: './PostEditor.component.html',
  styleUrls: ['./PostEditor.component.css']
})
export class PostEditorComponent implements OnInit { 
  postForm: FormGroup;

  constructor(private blogAPI: BlogAPI) {}

  ngOnInit() {
    this.postForm = new FormGroup ({
      title: new FormControl(),
      postMarkdown: new FormControl(),
    });
  }

  onSubmit() {
    const formModel = this.postForm.value;
    let post: Post = {
        Title: <string>formModel.title,
        Topic: 0,
        Body: <string>formModel.postMarkdown,
        Format: PostFormat.PostFormatMarkdown
      };
    this.blogAPI.PostPost(post).subscribe(
      () => {console.log("post saved")},
      (err) => {console.log(err.message)}
    );
  }

}
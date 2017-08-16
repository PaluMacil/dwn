import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'post-editor',
  templateUrl: './PostEditor.component.html',
  styleUrls: ['./PostEditor.component.css']
})
export class PostEditorComponent implements OnInit { 
  postMarkdown: string;

  constructor() {}

  ngOnInit() {
  }

  post = new FormGroup ({
    topic: new FormControl(),
    postMarkdown: new FormControl()
  });
}
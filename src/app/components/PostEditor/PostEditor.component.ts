import { Component, OnInit } from '@angular/core';

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
}
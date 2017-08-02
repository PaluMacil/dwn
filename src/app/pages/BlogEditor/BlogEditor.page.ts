import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'blog-editor-page',
  templateUrl: './BlogEditor.page.html',
  styleUrls: ['./BlogEditor.page.css']
})
export class BlogEditorPage implements OnInit { 
  constructor() {}

  ngOnInit() {
  }
}

export enum EditorMode {
  New,
  Edit
}
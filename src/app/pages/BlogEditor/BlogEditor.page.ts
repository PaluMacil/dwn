import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'blog-editor-page',
  templateUrl: './BlogEditor.page.html',
  styleUrls: ['./BlogEditor.page.css']
})
export class BlogEditorPage implements OnInit { 
  mode: EditorMode;
  EditorMode = EditorMode;

  constructor(route: ActivatedRoute) {
    this.mode = route.snapshot.data.mode;
    console.log(route);
  }

  ngOnInit() {
  }
}

export enum EditorMode {
  New,
  Edit
}
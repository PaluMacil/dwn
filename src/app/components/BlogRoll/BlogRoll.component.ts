import { Component, OnInit, Input } from '@angular/core';
import { BlogMode } from '../../services/BlogAPI.service';

@Component({
  selector: 'blog-roll',
  templateUrl: './BlogRoll.component.html',
  styleUrls: ['./BlogRoll.component.css']
})
export class BlogRollComponent implements OnInit { 
  @Input('mode') mode: BlogMode;

  constructor() {}

  ngOnInit() {
  }
}
import { Component, OnInit, Input } from '@angular/core';
import { MdRenderService } from '@nvxme/ngx-md-render';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { Post } from '../../services/BlogAPI.service';

@Component({
  selector: 'post',
  templateUrl: './Post.component.html',
  styleUrls: ['./Post.component.css']
})
export class PostComponent implements OnInit { 
  @Input('source') source: Post;
  rendered: SafeHtml;

  constructor(private sanitizer: DomSanitizer,
              private mdRender: MdRenderService) {}

  ngOnInit() {
    const html = this.mdRender.render(this.source.Body);
    this.rendered = this.sanitizer.bypassSecurityTrustHtml(html);
  }
}
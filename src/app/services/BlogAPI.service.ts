import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { UserInfo } from './AccountAPI.service'

@Injectable()
export class BlogAPI {
  constructor(private http: HttpClient) {}

  public GetPost(id: number): Observable<Post> {
    let params = new HttpParams()
      .set('ID', String(id));
    return this.http.get<Post>('/api/blog/post', {params});
  }

  public GetBlogRoll(): Observable<BlogRoll> {
    return this.http.get<BlogRoll>('/api/blog/roll');
  }
}

export interface BlogRoll {
  Posts:          Array<Post>;
  StartIndex:     number;
  PageSize:       number;
  TotalAvailable: number;
}

export interface Post {
  ID:       number;
  Title:    string;
  Slug:     string;
  Author:   UserInfo;
  Body:     string;
  Format:   PostFormat;
  IsStub:   boolean;
  Modified: string;
  Created:  string;
}

export const enum PostFormat {
    PostFormatPlain = 1,
    PostFormatMarkdown
}
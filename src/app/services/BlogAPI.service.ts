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

  public PostPost(post: Post): Observable<void> {
    return this.http.post<void>('/api/blog/post', post);
  }

  public PutPost(post: Post): Observable<void> {
    return this.http.put<void>('/api/blog/post', post);
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
  Topic:    Topic;
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
    PostFormatPlain = 0,
    PostFormatMarkdown
}

export const enum Topic {
    TopicAll = 0,
    TopicPersonal,
    TopicTech,
    TopicFood,
    TopicFun
}
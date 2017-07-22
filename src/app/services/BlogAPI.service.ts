import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { UserInfo } from './AccountAPI.service'

@Injectable()
export class BlogAPI {
  constructor(private http: HttpClient) {}

  public GetPost(id: number): Observable<Post> {
    return this.http.get<Post>(`/api/blog/post/${id}`)
  }
}

export interface BlogRoll {
  Posts:          Array<Post>;
  StartIndex:     number
  TotalAvailable: number
}

export interface Post {
  ID:     number;
  Title:  string;
  Author: UserInfo;
  Body:   string;
  IsStub: boolean;
}
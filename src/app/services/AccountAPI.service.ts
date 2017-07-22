import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class AccountAPI {
  constructor(private http: HttpClient) {}

  public Token(email: string, password: string): Observable<Session> {
    let params = new HttpParams()
      .set('email', email)
      .set('password', password);

    return this.http.post<Session>('/api/account/token', {params})
  }
}

export interface Session {
			Token:     string;
			User:      User;
			CreatedAt: string;
			HeartBeat: string;
    }
    
export interface User {
	ID:        number;
	Role:      number;
	Email:     string;
	Name:      string;
	CreatedAt: string;
}

export interface UserInfo {
  ID:      number;
  Name:    string;
  IsAdmin: boolean;
}
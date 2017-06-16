import { DwnPage } from './app.po';

describe('dwn App', () => {
  let page: DwnPage;

  beforeEach(() => {
    page = new DwnPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!!');
  });
});

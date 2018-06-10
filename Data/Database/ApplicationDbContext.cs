using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using Dwn.Data.Models.Identity;
using Dwn.Data.Models.Blogging;
using Dwn.Data.Models.Authoring;
using Dwn.Data.Models.Pages;

namespace Dwn.Data.Database
{
    public class ApplicationDbContext : IdentityDbContext<ApplicationUser, ApplicationRole, Guid>
    {
        public DbSet<AuthoredContent> AuthoredContent { get; set; }
        public DbSet<BlogCategory> BlogCategories { get; set; }
        public DbSet<Post> Posts { get; set; }
        public DbSet<Page> Pages { get; set; }
        public DbSet<Layout> Layouts { get; set; }

        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        { }

        protected override void OnModelCreating(ModelBuilder builder)
        {
            base.OnModelCreating(builder);
            
            builder.Entity<AuthoredContent>()
                .HasOne(ac => ac.Author)
                .WithMany(u => u.AuthoredContent)
                .OnDelete(DeleteBehavior.SetNull);
        }
    }
}
